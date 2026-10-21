<?php
echo 'before';
class MissingContract { public function make(): MissingBase { return new MissingBase; } }
class MissingImplementation extends MissingContract { public function make(): MissingDerived { return new MissingDerived; } }
echo 'after';
