#include <iostream>

template <typename T>
struct StackNode {
    T data;
    StackNode<T> *next;
};

template <typename T>
class LinkStack {
public:
    LinkStack();
    ~LinkStack();
public:
    bool Push(const T &e);
    bool Pop(T& e);
    bool GetTop(T &e);
    void DispList();
    int ListLength();
    bool Empty();

private:
    StackNode<T> *m_top;
    int m_length;
};

template <typename T>
LinkStack<T>::LinkStack() {
    m_top = nullptr;
    m_length = 0;
}

template <typename T>
bool LinkStack<T>::Push(const T &e) {
    StackNode<T> *node = new StackNode<T>;
    node->data = e;
    node->next = m_top;
    m_top = node;
    m_length++;
    return true;
}

template <typename T>
bool LinkStack<T>::Pop(T &e) {
    if (Empty() == true) {
        return false;
    }
    StackNode<T> *p_willdel = m_top;
    m_top = m_top->next;
    m_length--;
    e = p_willdel->data;
    delete p_willdel;
    return true;
}

template <typename T>
bool LinkStack<T>::GetTop(T &e) {
    if (Empty() == true) {
        return false;
    }
    e = m_top->data;
    return true;
}

template <class T>
void LinkStack<T>::DispList() {
    if (Empty() == true) {
        return;
    }

    StackNode<T> *p = m_top;
    while (p != nullptr) {
        std::cout << p->data << " ";
        p = p->next;
    }
    std::cout << std::endl;
}

template <class T>
int LinkStack<T>::ListLength() {
    return m_length;
}

template <class T>
bool LinkStack<T>::Empty() {
    if (m_top == nullptr) {
        return true;
    }
    return false;
}

template <typename T>
LinkStack<T>::~LinkStack() {
    T tmpnousevalue = {0};
    while (Pop(tmpnousevalue) == true) {}
}

int main(void) {
    LinkStack<int> slinkobj;
    slinkobj.Push(12);
    slinkobj.Push(24);
    slinkobj.Push(48);
    slinkobj.Push(100);

    int eval3 = 0;
    slinkobj.Pop(eval3);
    slinkobj.DispList();

    return 0;
}