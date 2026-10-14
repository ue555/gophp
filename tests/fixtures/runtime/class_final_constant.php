<?php
class ConstantParent { final public const VALUE = 1; }
echo "before\n";
class ConstantChild extends ConstantParent { public const VALUE = 2; }
echo "after\n";
