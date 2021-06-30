#!/usr/bin/env python3

import requests
import subprocess

class Updater:
    def __init__(self, name, version_file):
        self.name = name
        self.version_file = version_file

    def update(self, latest_version):
        updated = self.write_version_file(latest_version)
        if updated:
            print(" => Committing changes to git")
            self.commit_changes(latest_version)

    def write_version_file(self, latest_version):
        with open(self.version_file, "r+") as f:
            current = f.read()
            print(f" -> Current version = {current} (from {self.version_file})")

            f.seek(0)
            f.write(latest_version)
            f.truncate()

            return current != latest_version

    def commit_changes(self, latest_version):
        subprocess.run(["git", "add", self.version_file])
        subprocess.run(["git", "commit", "-m", f"Update {self.name} to version '{latest_version}'"])

class NpmUpdater(Updater):
    def __init__(self, name, version_file, npm_package):
        self.npm_package = npm_package;
        super().__init__(name, version_file)

    def update(self):
        print(f"Updating {self.name} via npm")

        latest_version = self.fetch_from_registry()
        super().update(latest_version)

    def fetch_from_registry(self):
        url = f"https://registry.npmjs.org/{self.npm_package}"
        res = requests.get(url)
        data = res.json()
        latest = data["dist-tags"]["latest"]

        print(f" -> Latest version = {latest} (from {url})")
        return latest

def main():
    projects = [
        NpmUpdater(
            name="MeshCentral",
            version_file="meshcentral/VERSION",
            npm_package="meshcentral",
        ),
    ]

    for project in projects:
        project.update()

if __name__ == "__main__":
    main()
