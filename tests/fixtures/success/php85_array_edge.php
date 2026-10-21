<?php
if (PHP_VERSION_ID >= 80500) {
    $values = [9 => 'first', 2 => 'last'];
    echo function_exists('array_first') ? 'yes|' : 'no|';
    echo array_first(array: $values), '|', array_last(array: $values), '|';
    echo array_first([]) === null && array_last([]) === null ? 'empty' : 'bad';
} else {
    echo function_exists('array_first') ? 'yes' : 'no';
}
