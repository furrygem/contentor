#!/usr/bin/env python3

import subprocess
import argparse
import sys
import requests
import etcd3


def start():
    subprocess.call(["./build/app"], shell=True)


def generate(_):
    print('generate')
    subprocess.run(["openssl", "genpkey", "-algorithm", "ed25519", "-out",
                    "key.pem"])
    subprocess.run(["openssl", "pkey", "-in", "key.pem", "-pubout", "-out",
                    "key.pub"])


def use_private_key(args):
    print("use_private_key")
    subprocess.run(["openssl", "pkey", "-in", args.f, "-pubout", "-out",
                    "key.pub"])


def use_http(args):
    print("use_http")
    r = requests.get(args.u)
    if r.status_code != 200:
        raise RuntimeError(f"Server at {args.u} returned {r.status_code}")
    open("key.pub", "wb").write(r.raw)


def use_etcd(args):
    print("use etcd")
    c = etcd3.client(args.host, args.port)
    val = c.get(args.k)
    if not val:
        sys.exit(255)
    open("key.pub", "wb").write(val[0])


if __name__ == "__main__":
    parser = argparse.ArgumentParser()

    subparser = parser.add_subparsers(required=True)
    generate_subparser = subparser.add_parser("generate")
    generate_subparser.set_defaults(func=generate)

    use_private_key_subparser = subparser.add_parser("use-private-key")
    use_private_key_subparser.set_defaults(func=use_private_key)
    use_private_key_subparser.add_argument("-f", action="store",
                                           help="Private key file",
                                           required=True)

    use_http_subparser = subparser.add_parser("use-http")
    use_http_subparser.set_defaults(func=use_http)
    use_http_subparser.add_argument('-u', action="store", help="URL")

    use_etcd_subparser = subparser.add_parser("use-etcd")
    use_etcd_subparser.set_defaults(func=use_etcd)
    use_etcd_subparser.add_argument("-host",
                                    help="etcd host",
                                    default="localhost")
    use_etcd_subparser.add_argument("-port",
                                    help="etcd port",
                                    default=2379,
                                    type=int)

    use_etcd_subparser.add_argument("-k",
                                    help="etcd key",
                                    default="contentor/public-key.pem")

    n = parser.parse_args()
    n.func(n)
    start()
