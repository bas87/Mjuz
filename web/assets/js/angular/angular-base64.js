(function() {
    'use strict';

    /*
     * Encapsulation of Nick Galbreath's base64.js library for AngularJS
     * Original notice included below
     */

    /*
     * Copyright (c) 2010 Nick Galbreath
     * http://code.google.com/p/stringencoders/source/browse/#svn/trunk/javascript
     *
     * Permission is hereby granted, free of charge, to any person
     * obtaining a copy of this software and associated documentation
     * files (the "Software"), to deal in the Software without
     * restriction, including without limitation the rights to use,
     * copy, modify, merge, publish, distribute, sublicense, and/or sell
     * copies of the Software, and to permit persons to whom the
     * Software is furnished to do so, subject to the following
     * conditions:
     *
     * The above copyright notice and this permission notice shall be
     * included in all copies or substantial portions of the Software.
     *
     * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
     * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
     * OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
     * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
     * HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
     * WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
     * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
     * OTHER DEALINGS IN THE SOFTWARE.
     */

    /* base64 encode/decode compatible with window.btoa/atob
     *
     * window.atob/btoa is a Firefox extension to convert binary data (the "b")
     * to base64 (ascii, the "a").
     *
     * It is also found in Safari and Chrome.  It is not available in IE.
     *
     * if (!window.btoa) window.btoa = base64.encode
     * if (!window.atob) window.atob = base64.decode
     *
     * The original spec's for atob/btoa are a bit lacking
     * https://developer.mozilla.org/en/DOM/window.atob
     * https://developer.mozilla.org/en/DOM/window.btoa
     *
     * window.btoa and base64.encode takes a string where charCodeAt is [0,255]
     * If any character is not [0,255], then an exception is thrown.
     *
     * window.atob and base64.decode take a base64-encoded string
     * If the input length is not a multiple of 4, or contains invalid characters
     *   then an exception is thrown.
     */

    angular.module('base64', []).constant('$base64', (function() {

        var _keyStr = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=";

        function encode(e) {
            var t = "";
            var n, r, i, s, o, u, a;
            var f = 0;
            e = _utf8_encode(e);
            while (f < e.length) {
                n = e.charCodeAt(f++);
                r = e.charCodeAt(f++);
                i = e.charCodeAt(f++);
                s = n >> 2;
                o = (n & 3) << 4 | r >> 4;
                u = (r & 15) << 2 | i >> 6;
                a = i & 63;
                if (isNaN(r)) {
                    u = a = 64
                } else if (isNaN(i)) {
                    a = 64
                }
                t = t + _keyStr.charAt(s) + _keyStr.charAt(o) + _keyStr.charAt(u) + _keyStr.charAt(a)
            }
            return t
        }

        function decode(e) {
            var t = "";
            var n, r, i;
            var s, o, u, a;
            var f = 0;
            e = e.replace(/[^A-Za-z0-9\+\/\=]/g, "");
            while (f < e.length) {
                s = _keyStr.indexOf(e.charAt(f++));
                o = _keyStr.indexOf(e.charAt(f++));
                u = _keyStr.indexOf(e.charAt(f++));
                a = _keyStr.indexOf(e.charAt(f++));
                n = s << 2 | o >> 4;
                r = (o & 15) << 4 | u >> 2;
                i = (u & 3) << 6 | a;
                t = t + String.fromCharCode(n);
                if (u != 64) {
                    t = t + String.fromCharCode(r)
                }
                if (a != 64) {
                    t = t + String.fromCharCode(i)
                }
            }
            t = _utf8_decode(t);
            return t
        }

        function _utf8_encode(e) {
            e = e.replace(/\r\n/g, "\n");
            var t = "";
            for (var n = 0; n < e.length; n++) {
                var r = e.charCodeAt(n);
                if (r < 128) {
                    t += String.fromCharCode(r)
                } else if (r > 127 && r < 2048) {
                    t += String.fromCharCode(r >> 6 | 192);
                    t += String.fromCharCode(r & 63 | 128)
                } else {
                    t += String.fromCharCode(r >> 12 | 224);
                    t += String.fromCharCode(r >> 6 & 63 | 128);
                    t += String.fromCharCode(r & 63 | 128)
                }
            }
            return t
        }

        function _utf8_decode(e) {
            var t = "";
            var n = 0;
            var r = c1 = c2 = 0;
            while (n < e.length) {
                r = e.charCodeAt(n);
                if (r < 128) {
                    t += String.fromCharCode(r);
                    n++
                } else if (r > 191 && r < 224) {
                    c2 = e.charCodeAt(n + 1);
                    t += String.fromCharCode((r & 31) << 6 | c2 & 63);
                    n += 2
                } else {
                    c2 = e.charCodeAt(n + 1);
                    c3 = e.charCodeAt(n + 2);
                    t += String.fromCharCode((r & 15) << 12 | (c2 & 63) << 6 | c3 & 63);
                    n += 3
                }
            }
            return t
        }

        return {
            encode: encode,
            decode: decode
        };
    })());

})();
