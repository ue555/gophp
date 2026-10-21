<?php
class NameExceptionBox { public int $value = 1; }
function failName() { echo 'name|'; throw new Exception('stop'); }
function skippedValue() { echo 'rhs|'; return 2; }
$box = new NameExceptionBox;
try { $box->{failName()} = skippedValue(); } catch (Exception $e) { echo $e->getMessage(), '|'; }
echo $box->value, "\n";
