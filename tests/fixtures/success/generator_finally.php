<?php
function protectedSequence() {
    try {
        yield 1;
        yield 2;
    } finally {
        echo 'finally|';
    }
}
foreach (protectedSequence() as $value) {
    echo $value, '|';
}
