<?php
class TypedErrors { public int $count; public static int $shared; }
$object = new TypedErrors;
try { echo $object->count; } catch (Error $error) { echo 'instance'; }
echo '|';
try { echo TypedErrors::$shared; } catch (Error $error) { echo 'static'; }
echo '|';
try { $object->count = []; } catch (TypeError $error) { echo 'assignment'; }
echo '|', isset($object->count) ? 'set' : 'unset';
try { $object->count = 4; unset($object->count); echo $object->count; }
catch (Error $error) { echo '|unset'; }
