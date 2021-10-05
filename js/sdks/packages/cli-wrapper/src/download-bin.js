/*
 * Copyright 2021 by LunaSec (owned by Refinery Labs, Inc)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
const fs = require('fs')
const path = require('path')
const axios = require('axios')
const decompress = require('decompress');
const decompressTargz = require('decompress-targz');
const ProgressBar = require('progress')

const binUrl = require('./format-bin-url')
// const binUrl = 'https://cran.r-project.org/src/base/R-4/R-4.1.1.tar.gz' // A good endpoint to test downloading on

const binFolder = path.resolve(__dirname, '..', 'bin')
const tarPath = path.resolve(binFolder, 'cli_bin.tar.gz')

async function downloadBinary(url) {

  // make the bin folder if it doesn't exist
  !fs.existsSync(binFolder) && fs.mkdirSync(binFolder, { recursive: true })

  const {data, headers} = await axios({
    url,
    method: 'GET',
    responseType: 'stream'
  })

  const progressBar = new ProgressBar('-> Downloading LunaSec CLI bin [:bar] :percent :etas', {
    width: 40,
    complete: '=',
    incomplete: ' ',
    renderThrottle: 1,
    total: parseInt(headers['content-length'])
  })

  data.on('data', (chunk) => progressBar.tick(chunk.length))

  const writer = fs.createWriteStream(tarPath)
  data.pipe(writer)

  return new Promise((resolve, reject) => {
    writer.on('finish', resolve)
    writer.on('error', reject)
  })
}
async function downloadAndDecompressCli(){

  await downloadBinary(binUrl).catch((e) => {
    console.error('Error downloading LunaSec CLI bin from github: ', e)
  })
  // Delete old decompressed files here if needed
  console.log('Extracting...')
  await decompress(tarPath, path.resolve(__dirname, '..', 'bin'), {plugins: [decompressTargz()]}).catch((e) => {
    console.error('Error decompressing LunaSec CLI bin',e )
  })

  fs.unlinkSync(tarPath) // Deletes the tar file
}

downloadAndDecompressCli().then(() => {
  console.log('LunaSec CLI installed')
}).catch((e) => {
  console.error('Error Installing LunaSec CLI: ', e)
  process.exit(1)
})
