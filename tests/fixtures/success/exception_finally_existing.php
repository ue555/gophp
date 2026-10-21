<?php
$first = new Exception('first');
$second = new Exception('second', 0, $first);
try {
    try { throw new Error('active'); }
    finally { throw $second; }
} catch (Exception $caught) {
    echo $caught->getMessage(), '|';
    echo $caught->getPrevious()->getMessage(), '|';
    echo $caught->getPrevious()->getPrevious()->getMessage(), "\n";
}
