# Logger

This repo will store every changes on database on a signed blockchain

## CI/CD
To enable CI/CD, first, you have to protect this tags:
```
release-stag-*
release-prod-*
```

After protect, to deploy you just need to push a tag using this semantic:
```
release-<env>-<semver>

Ex.:
release-stag-0.0.1.rc.1
release-prod-0.0.1
```

## TODO
- [ ] Listen Pub/Sub and create new blocks based on a message
- [ ] Make get all blocks work in batches