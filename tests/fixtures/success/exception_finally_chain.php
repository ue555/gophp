<?php
try {
    try { throw new Exception('first'); }
    finally { throw new Error('second'); }
} catch (Error $caught) {
    echo $caught->getMessage(), '|', $caught->getPrevious()->getMessage(), '|';
    echo ($caught->getPrevious()->getPrevious() === null ? 'end' : 'wrong'), "\n";
}
