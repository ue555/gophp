<?php
namespace SpecialTypes;
class Contract { public function make(): Root { return new Root; } }
class Implementation extends Contract { public function make(): Leaf { return new Leaf; } }
class Leaf extends Root {}
class Root {}
echo (new Implementation)->make() instanceof Root ? 'namespace' : 'wrong';
