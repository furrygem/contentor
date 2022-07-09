#!/usr/bin/env python3

import argparse
import subprocess
import sys
import etcd3


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("-m", help="Run alembic migrations on startup",
                        action="store_true")
    approach_group = parser.add_mutually_exclusive_group()

    approach_group.add_argument("--generate", help="Generate private key and overwrite\
        the existing one (if exists)", action="store_true")

    etcd_group = parser.add_argument_group("etcd")
    # group.add_argument("--use-etcd")
    etcd_group.add_argument("--use-etcd", action="store_true",
                            help="Use etcd to upload public key")
    etcd_group.add_argument("--etcd-host", default="localhost")
    etcd_group.add_argument("--etcd-port", type=int, default=2379)
    etcd_group.add_argument("--etcd-key", default="contentor/public-key.pem")

    parser.add_argument("-u", action="append", default=[])

    args = parser.parse_args()
    print(args)
    if args.generate:
        subprocess.call(["openssl", "genpkey", "-algorithm", "ed25519", "-out",
                         "key.pem"])

    subprocess.call(["openssl", "pkey", "-in", "key.pem", "-pubout", "-out",
                    "key.pub"])

    if args.use_etcd:
        etcd_client = etcd3.client(args.etcd_host, args.etcd_port)
        etcd_client.put(args.etcd_key, open("key.pub", 'rb').read())
        print(etcd_client.get(args.etcd_key))
        etcd_client.close()

    if args.m:
        if status := subprocess.call(["alembic", "upgrade", "head"]) != 0:
            sys.exit(status)
    subprocess.call(["uvicorn", "main:app", *args.u])
