<?php
enum Direction { case North; case South; public const Preferred = self::North; }
enum Level: int { case Low = 1; case High = 2; }
enum Code: string { case Ready = 'r'; case Done = 'd'; }
$cases = Direction::cases();
echo count($cases), '|', $cases[0]->name, '|';
echo Direction::North === $cases[0] ? 'same|' : 'different|';
echo Direction::Preferred === Direction::North ? 'alias|' : 'bad|';
echo Direction::North === Direction::South ? 'bad|' : 'distinct|';
echo Direction::North == Direction::South ? 'bad|' : 'looseDistinct|';
echo Level::from(2)->name, '|', Level::tryFrom(8) === null ? 'missing|' : 'bad|';
echo Code::from(value: 'r')->value;
