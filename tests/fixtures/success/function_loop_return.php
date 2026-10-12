<?php
function findValue() {
    for ($i = 0; $i < 3; ++$i) {
        do { while (true) { return $i + 7; } } while (1 / 0);
    }
    echo 'bad';
}
function early() { for (; true; 1 / 0) { return 4; } }
echo findValue(), ':', early(), "\n";
