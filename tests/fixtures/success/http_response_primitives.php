<?php
echo strrchr('a/b/c', '/'), '|', strrchr('a/b/c', '/', true), '|';
echo strrchr('abc', 'x') === false ? 'missing' : 'bad';
echo '|', ini_set('default_mimetype', ''), '|', ini_get('default_mimetype') === '' ? 'empty' : 'bad';
