<?php
class HttpAliasSource {}
class_alias(HttpAliasSource::class, 'HttpAliasTarget');
$data = new ArrayObject(['z' => 1, 'a' => 2]);
$previous = $data->exchangeArray(['z' => 3, 'a' => 4]);
$data->ksort();
$keys = array_keys($data->getArrayCopy());
$sorted = ['z' => 1, 'a' => 2];
uksort($sorted, static fn ($a, $b) => strcmp($a, $b));
echo new HttpAliasTarget instanceof HttpAliasSource ? 'alias' : 'bad';
echo '|', $previous['z'], '|', implode(',', $keys), '|', implode(',', array_keys($sorted));
echo '|', strrev('abc'), '|', bin2hex('A');
