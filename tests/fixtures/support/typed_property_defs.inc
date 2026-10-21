<?php
class CrossFileTyped { public int $count = 1; }
function weakTypedWrite($object) { $object->count = '7'; }
