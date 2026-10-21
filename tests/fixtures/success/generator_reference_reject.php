<?php
function byValue() {
    yield 1;
}
try {
    foreach (byValue() as &$value) {
        echo 'unreachable';
    }
} catch (Exception $e) {
    echo 'caught';
}
