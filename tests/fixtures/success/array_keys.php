<?php
$a = ['1' => 'first', 1 => 'int', true => 'bool', 1.0 => 'float',
      '01' => 'leading', '+1' => 'plus', ' 1' => 'space', '-0' => 'negative-zero',
      false => 'zero', '' => 'empty'];
echo $a[1], '|', $a['01'], '|', $a['+1'], '|', $a[' 1'], '|', $a['-0'], '|', $a[0], '|', $a[''], "\n";
echo ($a === [1 => 'float', '01' => 'leading', '+1' => 'plus', ' 1' => 'space', '-0' => 'negative-zero', 0 => 'zero', '' => 'empty']) + 0, "\n";
