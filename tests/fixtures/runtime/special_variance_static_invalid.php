<?php
echo 'before';
class StaticContract { public function make(): static { return new static; } }
class StaticImplementation extends StaticContract { public function make(): self { return new self; } }
echo 'after';
