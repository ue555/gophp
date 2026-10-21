<?php
try {
    try { throw new Exception('prior'); }
    finally { strlen([]); }
} catch (TypeError $error) {
    echo $error->getPrevious()->getMessage(), '|', count($error->getTrace()), "\n";
}
