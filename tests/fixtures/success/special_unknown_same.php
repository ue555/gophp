<?php
class UnknownContract { public function pass(MissingType $value): MissingType { return $value; } }
class UnknownImplementation extends UnknownContract { public function pass(MissingType $value): MissingType { return $value; } }
if (false) { class MissingType {} }
echo 'same';
