<?php
class LateBase { public function make(): static { return new self; } }
class LateChild extends LateBase {}
try { (new LateChild)->make(); echo 'wrong'; } catch (TypeError $error) { echo 'caught'; }
