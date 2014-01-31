# 1mage

`1mage` is an image upload bay, written in Node.js. To run `1mage`, simply start a `mongod` process in the background, then compile and run `1mage`:

    make
    node ./build/server/server.js <port> <path>

Where `<port>` is the port to use and `<path>` is the path to a directory which the server will use to store image files.

### License

Copyright (c) 2014 Alex Nichol & Jon Loeb

This is under the [GPL](http://www.gnu.org/licenses/gpl.html).
