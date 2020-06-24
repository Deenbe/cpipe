#!/usr/bin/env node

for (let i = 0; i < 100; i++) {
    console.log(JSON.stringify({
        pk: `${i}`,
        value: Math.random() * 10
    }));
}
