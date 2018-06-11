#!/usr/bin/env python3

import subprocess
import base64
import sys

GO_OS_ARCH_LIST = [
    ["darwin", "amd64"],
    ["linux", "386"],
    ["linux", "amd64"],
    ["linux", "arm"],
    ["linux", "arm64"],
    ["linux", "mips", "softfloat"],
    ["linux", "mips", "hardfloat"],
    ["linux", "mipsle"],
    ["linux", "mips64"],
    ["linux", "mips64le"],
    ["windows", "386"],
    ["windows", "amd64"]
              ]


def go_build_zip():
    subprocess.check_call("go get -v github.com/wheelerlaw/octodns/main", shell=True)
    for o, a, *p in GO_OS_ARCH_LIST:
        zip_name = "octodns-" + o + "-" + a + ("-" + (p[0] if p else "") if p else "")
        binary_name = zip_name + (".exe" if o == "windows" else "")
        version = subprocess.check_output("git describe --tags", shell=True).decode()
        mipsflag = (" GOMIPS=" + (p[0] if p else "") if p else "")
        try:
            subprocess.check_call("GOOS=" + o + " GOARCH=" + a + mipsflag + " CGO_ENABLED=0" + " go build -ldflags \"-s -w " +
                                  "-X main.version=" + version + "\" -o " + binary_name + " main/main.go", shell=True)
            subprocess.check_call("zip " + zip_name + ".zip " + binary_name + " " + IP_NETWORK_SAMPLE_DICT["name"] + " " +
                                  DOMAIN_SAMPLE_DICT["name"] + " hosts_sample domain_white_sample config.json", shell=True)
        except subprocess.CalledProcessError:
            print(o + " " + a + " " + (p[0] if p else "") + " failed.")

if __name__ == "__main__":
    if "build" in sys.argv:
        go_build_zip()
