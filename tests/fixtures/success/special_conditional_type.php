<?php
if (true) { class ConditionalRoot {} class ConditionalLeaf extends ConditionalRoot {} }
class ConditionalContract { public function make(): ConditionalRoot { return new ConditionalRoot; } }
class ConditionalImplementation extends ConditionalContract { public function make(): ConditionalLeaf { return new ConditionalLeaf; } }
echo (new ConditionalImplementation)->make() instanceof ConditionalRoot ? 'conditional' : 'wrong';
