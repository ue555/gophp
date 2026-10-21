<?php
foreach ([240, 241, 256, 1024, 1025] as $size) {
    $input = str_repeat('x', $size);
    $hash = hash('xxh3', $input);
    echo strlen($hash), ':';
    if ($size > 240) {
        echo $hash === substr(hash('xxh128', $input), 16) ? 'same' : 'different';
    } else {
        echo 'short';
    }
    echo '|';
}
