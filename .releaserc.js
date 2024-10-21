module.exports = {
    "branches": [
        {
            "name": "develop",
            "prerelease": true
        },
        {
            "name": "main",
            "prerelease": false
        }
    ],
    "plugins": [
        "@semantic-release/commit-analyzer",
        "@semantic-release/release-notes-generator",
        [
            "@semantic-release/git",
            {
                "assets": [
                    "CHANGELOG.md"
                ]
            }
        ],
        "@semantic-release/github"
    ],
    "repositoryUrl": "https://github.com/nilpntr/certmetrics-exporter.git",
}