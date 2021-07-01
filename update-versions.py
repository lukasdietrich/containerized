#!/usr/bin/env python3

import requests
import subprocess

class Updater:

    def __init__(self, name, version_file, version_fetcher):
        self.name = name
        self.version_file = version_file
        self.version_fetcher = version_fetcher

    def update(self):
        print("=> Updating {}".format(self.name))

        latest_version = self.version_fetcher.fetch_latest()
        current_version = self.write_version_file(latest_version)

        print("-> Current = {}, Latest = {}".format(current_version, latest_version))

        if latest_version != current_version:
            print("=> Committing changes to git")
            self.commit_changes(latest_version)

    def write_version_file(self, latest_version):
        with open(self.version_file, "r+") as f:
            current = f.read()

            f.seek(0)
            f.write(latest_version)
            f.truncate()

            return current

    def commit_changes(self, latest_version):
        message = "Update {} to version {}".format(self.name, latest_version)

        self.execute_git(["add", self.version_file])
        self.execute_git(["commit", "-m", message])

    def execute_git(self, args):
        command = ["git"] + args
        print("-> Executing {}".format(command))
        subprocess.run(command)

class NpmFetcher:

    def __init__(self, npm_package):
        self.npm_package = npm_package;

    def fetch_latest(self):
        url = "https://registry.npmjs.org/{}".format(self.npm_package)
        print("-> Fetching {}".format(url))

        res = requests.get(url)
        data = res.json()
        latest = data["dist-tags"]["latest"]

        return latest

def main():
    projects = [
        Updater(
            name="MeshCentral",
            version_file="meshcentral/VERSION",
            version_fetcher=NpmFetcher("meshcentral"),
        ),
    ]

    for project in projects:
        project.update()

if __name__ == "__main__":
    main()
