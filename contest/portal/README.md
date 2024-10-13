
`src/main.ts` から `dist/main.js` を生成。
```
$ npm run build
```

Terraform から `index.html` と `dist/main.js` を参照するので、 `dist/*` 配下もgitでバージョン管理することにする。
