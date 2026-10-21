<?php
try {
    try {
        try { throw new Exception('first'); }
        finally { throw new Error('second'); }
    } finally { throw new Exception('third'); }
} catch (Exception $error) {
    echo $error->getMessage(), '|';
    echo $error->getPrevious()->getMessage(), '|';
    echo $error->getPrevious()->getPrevious()->getMessage(), "\n";
}
