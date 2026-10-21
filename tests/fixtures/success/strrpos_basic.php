<?php
echo strrpos('a/b/c', '/'), '|';
echo strrpos('abc', 'z') === false ? 'missing' : 'wrong';
