export const THEME_OPTIONS = [
  { value: 'auto', label: 'Selon l’image' },
  { value: 'calm', label: 'Contemplatif' },
  { value: 'energetic', label: 'Joyeux' },
  { value: 'melancholy', label: 'Mélancolique' },
  { value: 'tense', label: 'Dramatique' },
  { value: 'neutral', label: 'Introspectif' },
] as const

export const THEME_LABELS: Record<string, string> = Object.fromEntries(THEME_OPTIONS.map(option => [option.value, option.label]))
