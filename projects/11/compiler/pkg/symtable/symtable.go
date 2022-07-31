package symtable

type SymbolTableEntryKind int

const (
	Static SymbolTableEntryKind = iota
	Field
	Arg
	Var
	None
)

type symbolTableEntry struct {
	name      string
	entryType string
	kind      SymbolTableEntryKind
	index     int
}

type SymbolTable struct {
	entries     map[string]symbolTableEntry
	staticIndex int
	fieldIndex  int
	argIndex    int
	varIndex    int
}

func New() *SymbolTable {
	return &SymbolTable{
		entries:     make(map[string]symbolTableEntry),
		staticIndex: 0,
		fieldIndex:  0,
		argIndex:    0,
		varIndex:    0,
	}
}

func (st *SymbolTable) Reset() {
	st.entries = make(map[string]symbolTableEntry)
	st.staticIndex = 0
	st.fieldIndex = 0
	st.argIndex = 0
	st.varIndex = 0
}

func (st *SymbolTable) Define(name string, entryType string, kind SymbolTableEntryKind) {
	var index int
	switch kind {
	case Static:
		index = st.staticIndex
		st.staticIndex++
	case Field:
		index = st.fieldIndex
		st.fieldIndex++
	case Arg:
		index = st.argIndex
		st.argIndex++
	case Var:
		index = st.varIndex
		st.varIndex++
	}

	st.entries[name] = symbolTableEntry{
		name:      name,
		entryType: entryType,
		kind:      kind,
		index:     index,
	}
}

func (st *SymbolTable) VarCount(kind SymbolTableEntryKind) int {
	switch kind {
	case Static:
		return st.staticIndex
	case Field:
		return st.fieldIndex
	case Arg:
		return st.argIndex
	case Var:
		return st.varIndex
	}
	return 0
}

func (st *SymbolTable) KindOf(name string) SymbolTableEntryKind {
	if entry, ok := st.entries[name]; ok {
		return entry.kind
	}
	return None
}

func (st *SymbolTable) TypeOf(name string) string {
	return st.entries[name].entryType
}

func (st *SymbolTable) IndexOf(name string) int {
	return st.entries[name].index
}
