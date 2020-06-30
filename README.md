# Eloquent ORM

Eloquent is a Golang ORM framework inspired by the famous PHP framework Laravel's Eloquent.

```yaml
package: models
imports:
- github.com/mylxsw/eloquent
meta:
  table_prefix: el_
models:
- name: user
  relations:
  - model: role
    rel: n-1
    foreign_key: role_id
    owner_key: id
    local_key: ""
    table: ""
    package: ""
    method: ""
  definition:
    table_name: user
    without_create_time: false
    without_update_time: false
    soft_delete: false
    fields:
    - name: id
      type: int64
      tag: json:"id"
    - name: name
      type: string
      tag: json:"name"
    - name: age
      type: int64
      tag: json:"age"
```

**This project is under heavy development !**

## Stargazers over time

[![Stargazers over time](https://starchart.cc/mylxsw/eloquent.svg)](https://starchart.cc/mylxsw/eloquent)