<?php
echo strtr('a\\b', '\\', '/'), '|';
echo strtr('one two', ['one' => '1', 'two' => '2']), '|', DIRECTORY_SEPARATOR;
