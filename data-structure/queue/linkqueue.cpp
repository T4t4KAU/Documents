#include <iostream>

template <typename T>
struct QueueNode {
    T data;
    QueueNode<T> *next;
};

template <typename T>
class LinkQueue {
public:
    LinkQueue();
    ~LinkQueue();
public:
    bool EnQueue(const T &e);
    bool DeQueue(T &e);
    bool GetHead(T &e);

    void DispList();
    int ListLength();
    bool IsEmpty();
private:
    QueueNode<T> *m_front;
    QueueNode<T> *m_rear;
    int m_length;
};

template <typename T>
LinkQueue<T>::LinkQueue() {
    m_front = new QueueNode<T>;
    m_front->next = nullptr;
    m_rear = m_front;
    m_length = 0;
}

template <typename T>
LinkQueue<T>::~LinkQueue() {
    QueueNode<T> *pnode = m_front->next;
    QueueNode<T> *ptmp;
    while (pnode != nullptr) {
        ptmp = pnode;
        pnode = pnode->next;
        delete ptmp;
    }
    delete m_front;
    m_front = m_rear = nullptr;
    m_length = 0;
}

template <typename T>
bool LinkQueue<T>::EnQueue(const T &e) {
    QueueNode<T> *node = new QueueNode<T>;
    node->data = e;
    node->next = nullptr;

    m_rear->next = node;
    m_rear = node;
    m_length++;

    return true;
}

template <typename T>
bool LinkQueue<T>::DeQueue(T &e) {
    if (IsEmpty() == true) {
        std::cout << "LinkQueue Empty" << std::endl;
        return false;
    }
    QueueNode<T> *p_willdel = m_front->next;
    e = p_willdel->data;

    m_front->next = p_willdel->next;
    if (m_rear == p_willdel) {
        m_rear = m_front;
    }

    delete p_willdel;
    m_length--;
    return true;
}

template <typename T>
bool LinkQueue<T>::GetHead(T &e) {
    if (IsEmpty() == true) {
        std::cout << "Link Queue Empty" <<std::endl;
        return false;
    }
    e = m_front->next->data;
    return true;
}

template <class T>
void LinkQueue<T>::DispList() {
    QueueNode<T> *p = m_front->next;
    while (p != nullptr) {
        std::cout << p->data << " ";
        p = p->next;
    }
    std::cout << std::endl;
}

template <class T>
int LinkQueue<T>::ListLength() {
    return m_length;
}

template <class T>
bool LinkQueue<T>::IsEmpty() {
    if (m_front == m_rear) {
        return true;
    }
    return false;
}

int main(void) {
    LinkQueue<int> lnobj;
    lnobj.EnQueue(150);

    int eval2 = 0;
    lnobj.DeQueue(eval2);
    lnobj.EnQueue(200);
    lnobj.EnQueue(700);
    lnobj.DispList();

    return 0;
}