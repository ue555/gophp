<?php
class NestedThrowTarget { public static function run() { echo 'run|'; } }
function methodName() { return 'run'; }
NestedThrowTarget::{true ? 'run' : throw new Exception('skip')}();
try {
    NestedThrowTarget::{methodName()}();
    NestedThrowTarget::{false ? 'run' : throw new Exception('stop')}();
} catch (Exception $e) {
    echo 'caught';
}
echo "\n";
