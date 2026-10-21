<?php
class ConstantWide { public const int|float VALUE = 1; }
class ConstantNarrow extends ConstantWide { public const float VALUE = 2; }
class ConstantIntToFloat { public const float VALUE = 3; }
echo ConstantNarrow::VALUE, '|', gettype(ConstantIntToFloat::VALUE);
