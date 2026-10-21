<?php
$path = tempnam(sys_get_temp_dir(), 'gophp-writable-');
echo (int)function_exists('is_writable'), (int)function_exists('is_writeable'), "\n";
echo (int)is_writable($path), (int)is_writeable(filename: $path), (int)is_writable(sys_get_temp_dir()), "\n";
echo (int)is_writable($path.'/missing'), (int)is_writable(''), (int)is_writable("a\0b"), "\n";
echo strlen(file_get_contents($path)), "\n";
unlink($path);
echo (int)is_writable($path), "\n";
