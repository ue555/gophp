<?php
function choose() {
    try {
        return 1;
    } finally {
        return 2;
    }
}
echo choose(), '|';
try {
    try { throw new Exception('nested'); }
    catch (Exception $e) { throw $e; }
} catch (Exception $e) {
    echo $e->getMessage(), "\n";
}
