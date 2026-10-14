<?php
try {
    require __DIR__ . '/../support/throwing.inc';
} catch (Exception $e) {
    echo $e->getMessage(), "\n";
}
