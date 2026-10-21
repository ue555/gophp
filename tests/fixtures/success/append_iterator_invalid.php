<?php
$combined = new AppendIterator;
try {
    $combined->append(new stdClass);
    echo 'bad';
} catch (TypeError $e) {
    echo 'rejected|';
}
echo count($combined->getArrayIterator());
