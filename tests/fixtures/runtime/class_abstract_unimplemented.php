<?php
echo "before\n";
abstract class AbstractMethodParent { abstract public function value(); }
class MissingMethodChild extends AbstractMethodParent {}
