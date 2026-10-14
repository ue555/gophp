<?php
$globalError = null;
function saveGlobal() {
    global $globalError;
    try { throw new Exception('global'); }
    catch (Exception $globalError) {}
}
function saveStatic() {
    static $staticError = null;
    try { throw new Exception('static'); }
    catch (Exception $staticError) {}
    echo gettype($staticError), '|', $staticError->getMessage(), "\n";
}
saveGlobal();
echo gettype($globalError), '|', $globalError->getMessage(), "\n";
saveStatic();
saveStatic();
