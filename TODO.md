# TODO for GCP Cloud Functions labeling Resources based on Audit Events
### Todo
- [ ] Implement proper Release Pipeline
- [ ] Might want to leverage "google.golang.org/genproto/googleapis/cloud/audit", so we can forget about types - Does PubSub sink with binary topic schema serialize using protobuf / [Publish messages of protobuf schema type](https://cloud.google.com/pubsub/docs/samples/pubsub-publish-proto-messages)?
- [ ] Test drive v2/EventArc
- [ ] Harden Default Compute Service Account (Revoke Editor)
- [ ] Stricter naming of terraform resources/modules
- [ ] Do org level instance actions role in `terraform`
- [ ] Error Handling/Retries (around actions function)

### In Progress
- [ ] CFT DevTools/E2E-Tests/Terratest

### Done ✓

- [x] Nothing