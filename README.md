Terraform Provider to generate Terraform State container Outputs only
==================

[![Build Status](https://travis-ci.org/sl1pm4t/terraform-provider-tfstate.svg?branch=master)](https://travis-ci.org/sl1pm4t/terraform-provider-tfstate)

A logical provider that can be used to generate a pseudo Terraform State file, containing `outputs` only.

**Why?**

Terraform Remote state is a useful way to transfer values between Terraform environments but it requires the state reader to have access to the entire state file, which may contain sensitive data.
Using this resources it's possible to generate a pseudo `.tfstate` file containing just the `outputs` without exposing internal details of the full Terraform config.
Additionally, permissions on the pseudo `.tfstate` file can be set independently of the real `.tfstate` file, or it could be stored in a different location that is more accessible to be consumed by downstream configs.

**Known Limitations**

Due to current limitations with the Terraform type system, it's only possible to use `string` typed values in the outputs.
Use Terraform [interpolation functions](https://www.terraform.io/docs/configuration/interpolation.html) such as `join`, `keys`, `values`, `list`, `zipmap` to encode/decode maps and lists to/from strings.

Using the provider
----------------------

**Basic Example**

```hcl
resource tfstate_outputs "test" {
  output {
    name  = "foo"
    value = "bar"
  }

  output {
    name  = "baz"
    value = "bam"
  }
}

resource "local_file" "state_outputs" {
  content  = "${tfstate_outputs.test.json}"
  filename = "${path.module}/terraform.tfout"
}

# In another module / config:

data "terraform_remote_state" "upstream" {
  backend = "local"

  config {
    path = "../terraform.tfout"
  }
}

output "upstream_foo" {
  value = "${data.terraform_remote_state.upstream.foo}"
}
```


Development Requirements
----------------------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Go](https://golang.org/doc/install) 1.10 (to build the provider plugin)

Building The Provider
-------------------------

Clone repository to: `$GOPATH/src/github.com/sl1pm4t/terraform-provider-tfstate`

```sh
$ mkdir -p $GOPATH/src/github.com/sl1pm4t; cd $GOPATH/src/github.com/sl1pm4t
$ git clone git@github.com:sl1pm4t/terraform-provider-tfstate
```
 
Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/sl1pm4t/terraform-provider-tfstate
$ make build
```




Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-tfstate
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.


```sh
$ make testacc
```
