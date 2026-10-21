<?php
class ForwardContract { public function make(): ForwardBase { return new ForwardBase; } }
class ForwardImplementation extends ForwardContract { public function make(): ForwardDerived { return new ForwardDerived; } }
class ForwardDerived extends ForwardBase {}
class ForwardBase {}
echo (new ForwardImplementation)->make() instanceof ForwardBase ? 'forward' : 'wrong';
