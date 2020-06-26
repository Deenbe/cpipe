# cpipe

> Pipe data streams to AWS services

![Build and Release](https://github.com/buddyspike/cpipe/workflows/Build%20and%20Release/badge.svg)

## Install cpipe

```sh
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Deenbe/cpipe/master/install.sh)"
```

## Usage

Save a single item to dynamodb.
```sh
echo '{"id": 1, "first_name": "barry", "last_name": "allen"}' | cpipe dynamodb --table dc_crowd
```

Similarly you could pipe the output of a program that generates a series of items to `cpipe` to store them in dynamodb.
```js
// gen.js
#!/usr/bin/env node

for (let i = 0; i < 100; i++) {
    console.log(JSON.stringify({
        id: `${i}`,
        value: Math.random() * 10
    }));
}
```

```sh
./gen.js | cpipe dynamodb --table sensor_data
```


