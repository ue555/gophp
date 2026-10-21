<?php
class TypedErrors { public int $value = 1; public int $missing; }
$box = new TypedErrors;
try { $alias =& $box->missing; } catch (Error $e) { echo 'uninitialized|'; }
$alias =& $box->value;
try { $alias = []; } catch (TypeError $e) { echo 'rejected|'; }
echo $box->value, '|';
$bad = [];
try { $box->value =& $bad; } catch (TypeError $e) { echo 'bind-rejected|'; }
echo $box->value, '|', gettype($bad);
