<?php
try {
    echo 'A';
    throw new Exception('bad', 7);
} catch (Error|Exception $e) {
    echo $e->getMessage(), ':', $e->getCode();
} finally {
    echo 'F', "\n";
}
